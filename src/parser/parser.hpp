#pragma once

#include <boost/program_options/variables_map.hpp>

#include <parser/node/master_node.hpp>
#include <parser/node/structure/parse_with_descriptor_node.hpp>

namespace po = boost::program_options;

namespace ql::parser {
    class Parser {
    private:
        using NodeFactory = std::function<std::shared_ptr<ParseWithDescriptorNode>(std::string&&, std::string_view const&, std::vector<std::string>&&,
                                                                                   AbstractNode::ParentRef)>;

        std::map<std::string, NodeFactory> m_NamesToNodes;

        template<typename TNode>
        void registerNode(std::string_view nodeName) {
            // TODO use forwarding?
            m_NamesToNodes.emplace(nodeName, [](std::string&& block, std::string_view const& body, std::vector<std::string>&& tokens,
                                                AbstractNode::ParentRef parent) {
                return std::make_shared<TNode>(std::move(block), body, std::move(tokens), parent);
            });
        }

        std::shared_ptr<AbstractNode> getNode(std::string const& nodeName,
                                              std::string&& blockWithInfo, std::string_view const& innerBlock, std::vector<std::string>&& tokens,
                                              AbstractNode::ParentRef parent);

        void recurseNodes(std::string_view code, std::weak_ptr<AbstractNode> const& parent, int depth = 0);

    public:
        Parser();

        std::shared_ptr<MasterNode> parse(po::variables_map& options);

        std::shared_ptr<MasterNode> getNodes(std::string code);
    };
}
