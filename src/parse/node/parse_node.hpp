#pragma once

#include <memory>
#include <vector>
#include <utility>

namespace ql::parse {
    class ParseNode {
    public:
        typedef std::vector<std::shared_ptr<ParseNode>> ChildrenRef;
        typedef std::weak_ptr<ParseNode> ParentRef;
    private:
        ChildrenRef m_Children;
        ParentRef m_Parent;
        std::string m_RawText;
    public:
        ParseNode(std::string&& rawText, ParentRef parent)
                : m_RawText(rawText), m_Parent(std::move(parent)) {}

        std::string_view getText() const { return m_RawText; }

        void addChild(const std::shared_ptr<ParseNode>& node);
    };
}
